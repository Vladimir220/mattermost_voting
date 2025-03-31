FROM tarantool/tarantool

COPY ./Db/Migrations/001_votings.lua /opt/tarantool/inittt.lua

WORKDIR /opt/tarantool

ENV TARANTOOL_USER=user
ENV TARANTOOL_PASS=qwerty

RUN ln -s /usr/share/zoneinfo/Europe/Moscow /etc/localtime
CMD ["tarantool", "/opt/tarantool/inittt.lua"]

EXPOSE 3301
