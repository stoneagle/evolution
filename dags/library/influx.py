from influxdb import InfluxDBClient, DataFrameClient
from library import conf


client = InfluxDBClient(
    host=conf.INFLUXDB_HOST,
    port=conf.INFLUXDB_PORT,
    username=conf.INFLUXDB_USER,
    password=conf.INFLUXDB_PASSWORD,
    database=conf.INFLUXDB_DBNAME,
)

dfclient = DataFrameClient(
    host=conf.INFLUXDB_HOST,
    port=conf.INFLUXDB_PORT,
    username=conf.INFLUXDB_USER,
    password=conf.INFLUXDB_PASSWORD,
    database=conf.INFLUXDB_DBNAME,
)


protocol = 'json'


def init_db():
    global client
    client.create_database(conf.INFLUXDB_DBNAME)
    return


def write_df(df, measurement):
    global dfclient
    dfclient.write_points(df, measurement, protocol=conf.INFLUXDB_PROTOCOL_JSON)
    return
