FROM alpine:3.17.0
ARG wekanAlerterDir
COPY --chmod=555 ./$wekanAlerterDir/wekan-alerter /app/wekan-alerter
COPY ./$wekanAlerterDir/mail-html.tmpl /app/mail-html.tmpl
#RUN chmod 555 /app/wekan-alerter
WORKDIR /app
CMD ["/app/datapi"]