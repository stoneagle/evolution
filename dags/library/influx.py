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


def write_df(df, measurement, tags):
    global dfclient
    if tags is None:
        dfclient.write_points(df, measurement, protocol=conf.INFLUXDB_PROTOCOL_JSON)
    else:
        dfclient.write_points(df, measurement, tags, protocol=conf.INFLUXDB_PROTOCOL_JSON)
    return


def reset_df(df, measurement, tags):
    global dfclient
    delete_measurement(measurement, tags)
    write_df(df, measurement, tags)
    return


def delete_measurement(measurement, tags):
    global dfclient
    dfclient.delete_series(measurement=measurement, tags=tags)
    return
