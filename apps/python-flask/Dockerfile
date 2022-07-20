FROM python:3.9.3-slim

COPY . .

RUN python -m venv venv
RUN venv/bin/python -m pip install --upgrade pip
RUN venv/bin/pip install -r requirements.txt
RUN venv/bin/pip install gunicorn

RUN chmod a+x boot.sh

ENTRYPOINT [ "./boot.sh" ]
