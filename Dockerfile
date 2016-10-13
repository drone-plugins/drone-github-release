FROM cdrx/fpm-centos:7

RUN yum install -y svn git

ADD drone-svn-release /bin/
ADD *.sh /bin/

ENTRYPOINT ["/bin/drone-svn-release"]
