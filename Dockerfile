FROM scratch

ADD ./dual_port /dual_port

EXPOSE 8080
EXPOSE 8443

ENTRYPOINT [ "/dual_port" ]