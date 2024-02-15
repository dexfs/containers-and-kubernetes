FROM scratch

WORKDIR /app

COPY api /app/

EXPOSE 8000

CMD [ "./api" ]