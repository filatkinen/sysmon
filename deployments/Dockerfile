FROM golang:1.21

RUN apt-get update
RUN apt-get install sudo
RUN apt-get install sysstat -y
RUN apt-get install tcpdump -y
RUN apt-get install net-tools -y
RUN apt-get install iproute2 -y

RUN adduser --disabled-password --gecos '' admin
RUN adduser admin sudo
RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers


WORKDIR /usr/src/app/

COPY . /usr/src/app/
RUN go build -v -o ./build/service  ./cmd/service

USER admin

ENTRYPOINT [ "sudo", "./build/service"]
CMD ["-config", "./configs/service.yaml"]
EXPOSE 50051
