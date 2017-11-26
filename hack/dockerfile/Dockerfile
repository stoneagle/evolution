FROM puckel/docker-airflow:1.8.2

USER root 

RUN apt-get -y update && \ 
apt-get -y install zlib1g-dev libhdf5-dev  libhdf5-serial-dev  wget

# 安装talib相关文件
RUN cd /tmp && \
    wget https://jaist.dl.sourceforge.net/project/ta-lib/ta-lib/0.4.0/ta-lib-0.4.0-src.tar.gz && \
    tar -zxvf ./ta-lib-0.4.0-src.tar.gz && \
    cd ta-lib && \
    ./configure --prefix=/usr && \
    make && \
    make install

COPY ./requirements.utils.txt /tmp/requirements.txt
RUN pip install --index https://pypi.mirrors.ustc.edu.cn/simple/ --user -r /tmp/requirements.txt

COPY ./requirements.tushare.txt /tmp/requirements.txt
RUN pip install --index https://pypi.mirrors.ustc.edu.cn/simple/ --user -r /tmp/requirements.txt

COPY ./hdf5 /usr/local/hdf5
RUN export HDF5_DIR=/usr/local/hdf5
ENTRYPOINT ["/bin/bash"]
