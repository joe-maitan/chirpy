FROM alpine

COPY chirpy /bin/chirpy

RUN apt-get update && apt-get install -y curl

ENV PORT
ENV DB_URL
ENV PLATFORM
ENV SECRET
ENV POLKA_API_KEY

HEALTHCHECK --interval=30s --timeout=10s --retries=3 --start-period=5s \
  CMD curl -f http://localhost:$PORT/api/healthz || exit 1

CMD ["/bin/chirpy"]
