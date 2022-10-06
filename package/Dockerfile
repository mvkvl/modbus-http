FROM golang:bullseye

# create systemd-enabled container to be able to test mbridge service

MAINTAINER "Mikhail Kantur” <mkantur@gmail.com>"
RUN apt update && apt install -y systemd && apt-get clean
CMD ["/bin/systemd"]
