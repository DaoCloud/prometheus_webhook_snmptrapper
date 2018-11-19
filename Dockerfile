FROM daocloud.io/daocloud/go-busybox:glibc

COPY prometheus_webhook_snmptrapper /
RUN chmod +x /prometheus_webhook_snmptrapper

COPY start.sh /
RUN chmod +x /start.sh
EXPOSE 5001
ENTRYPOINT [ "sh" ]
CMD ["/start.sh"]
