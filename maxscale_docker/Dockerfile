FROM centos:7

# -- maxscale --
ENV MAXSCALE_VERSION=2.3.6
RUN curl -sS https://downloads.mariadb.com/MariaDB/mariadb_repo_setup | bash -s -- --mariadb-maxscale-version="${MAXSCALE_VERSION}" && \
    yum -y install maxscale && \
    yum clean all



ENTRYPOINT ["/entrypoint.sh"]
