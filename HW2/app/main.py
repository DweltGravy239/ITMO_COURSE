# main.py
import json
import logging
import os
import random
import sys
import time

import httpx
from fastapi import FastAPI, Request, Response
from opentelemetry import trace
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor
from opentelemetry.sdk.resources import SERVICE_NAME, Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from prometheus_client import Counter, Histogram, generate_latest, CONTENT_TYPE_LATEST

resource = Resource(attributes={SERVICE_NAME: "api"})
provider = TracerProvider(resource=resource)
otlp_endpoint = os.environ.get("OTEL_EXPORTER_OTLP_ENDPOINT")
exporter = OTLPSpanExporter(endpoint=otlp_endpoint) if otlp_endpoint else OTLPSpanExporter()
provider.add_span_processor(BatchSpanProcessor(exporter))
trace.set_tracer_provider(provider)
tracer = trace.get_tracer("api")


class JsonFormatter(logging.Formatter):
    def format(self, record):
        ctx = trace.get_current_span().get_span_context()
        trace_id = format(ctx.trace_id, "032x") if ctx.is_valid else None
        return json.dumps({
            "timestamp": self.formatTime(record),
            "level": record.levelname,
            "message": record.getMessage(),
            "trace_id": trace_id,
        })


handler = logging.StreamHandler(sys.stdout)
handler.setFormatter(JsonFormatter())
logger = logging.getLogger("api")
logger.setLevel(logging.INFO)
logger.addHandler(handler)

REQUEST_COUNT = Counter("api_requests_total", "Total requests", ["endpoint"])
ERROR_COUNT = Counter("api_errors_total", "Total 5xx responses", ["endpoint"])
REQUEST_LATENCY = Histogram("api_request_duration_seconds", "Request latency", ["endpoint"], buckets=[0.1, 0.25, 0.5, 1, 1.5, 2, 2.5, 3, 4, 5])

app = FastAPI()
FastAPIInstrumentor.instrument_app(app)


@app.middleware("http")
async def metrics_middleware(request: Request, call_next):
    path = request.url.path
    if path == "/metrics":
        return await call_next(request)
    start = time.time()
    response = await call_next(request)
    REQUEST_LATENCY.labels(endpoint=path).observe(time.time() - start)
    REQUEST_COUNT.labels(endpoint=path).inc()
    if response.status_code >= 500:
        ERROR_COUNT.labels(endpoint=path).inc()
    return response


@app.get("/health")
def health():
    logger.info("health check ok")
    return "ok"


@app.get("/fail")
def fail():
    trace.get_current_span().set_status(trace.Status(trace.StatusCode.ERROR, "simulated failure"))
    logger.error("simulated failure triggered")
    return Response(content="internal error", status_code=500)


@app.get("/slow")
def slow():
    with tracer.start_as_current_span("slow-op"):
        delay = random.uniform(1, 3)
        time.sleep(delay)
        logger.info(f"slow operation finished in {delay:.2f}s")
    return {"slept": True}


@app.get("/load")
def load():
    logger.info("generating load")
    with httpx.Client(base_url="http://localhost:8000") as client:
        for _ in range(20):
            client.get("/health")
    return {"requests_sent": 20}


@app.get("/metrics")
def metrics():
    return Response(generate_latest(), media_type=CONTENT_TYPE_LATEST)
