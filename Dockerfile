FROM ubuntu:16.04
RUN apt-get update -y && \
    apt-get install -y python-pip python-dev

COPY ./requirements.txt /app/requirements.txt
WORKDIR /bidders
RUN pip install -r requirements.txt
COPY . /bidders
ENTRYPOINT [ "python" ]
CMD [ "bidder.py" ]
