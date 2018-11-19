FROM alpine:latest

RUN apk update
RUN apk add net-snmp

COPY PROMETHEUS-TRAPPER-MIB.txt /usr/share/snmp/mibs/

CMD ["/usr/sbin/snmptrapd", "-f", "-Lo", "-m", "PROMETHEUS-TRAPPER-MIB", "-M", "/usr/share/snmp/mibs"]