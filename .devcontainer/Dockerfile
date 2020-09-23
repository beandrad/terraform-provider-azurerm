ARG VARIANT="focal"

# To build without user input
ARG DEBIAN_FRONTEND=noninteractive

FROM ubuntu:${VARIANT}

ENV TERRAFORM_VERSION=0.12.29
ENV GO_VERSION=1.14.9

RUN apt-get update && apt-get install -y curl zip make gcc git

# Make shared download folder
RUN mkdir -p /tmp/docker-downloads

# Install Terraform
RUN curl -sSL https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip -o /tmp/docker-downloads/terraform.zip \
    && unzip /tmp/docker-downloads/terraform.zip \
    && mv terraform /usr/local/bin

# Install Go into /usr/local/go
RUN curl -fsS https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz -o /tmp/docker-downloads/golang.tar.gz \
    && tar -C /usr/local -xvf /tmp/docker-downloads/golang.tar.gz 

ENV PATH=$PATH:/usr/local/go/bin

# Remove shared download folder
RUN rm -rf /tmp/docker-downloads

ENV DEBIAN_FRONTEND=interactive
