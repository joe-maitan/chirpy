FROM alpine 

COPY chirpy .

ENV PORT=8991

CMD ["./chirpy"]