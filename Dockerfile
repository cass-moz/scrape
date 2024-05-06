FROM docker.io/debian:12-slim

RUN apt -y update && apt -y upgrade
RUN apt-get -y install \
    ca-certificates \
    curl \
    chromium \
    gnupg wget apt-transport-https
RUN mkdir -p /scrape/bin
COPY scrape-server scrape scrape-feed /scrape/bin/
RUN mkdir -p /scrape_data
VOLUME [ "/scrape_data" ]
EXPOSE 8080/tcp
# The default sqlite db will be in /scrape_data/scrape.db
ENTRYPOINT ["/scrape/bin/scrape-server"]