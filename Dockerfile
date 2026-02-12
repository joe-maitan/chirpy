FROM alpine

COPY chirpy .

ENV PORT=8991

HEALTHCHECK NONE

CMD ["./chirpy"]
