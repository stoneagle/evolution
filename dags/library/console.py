import sys


def write_tail(ctype):
    DATA_GETTING_TIPS = '\r\nFinish Getting data[' + ctype + ']\r\n'
    sys.stdout.write(DATA_GETTING_TIPS)
    sys.stdout.flush()


def write_head(ctype):
    DATA_GETTING_TIPS = 'Getting data[' + ctype + ']:\r\n'
    sys.stdout.write(DATA_GETTING_TIPS)
    sys.stdout.flush()


def write_create():
    sys.stdout.write('C')
    sys.stdout.flush()


def write_append():
    sys.stdout.write('A')
    sys.stdout.flush()


def write_number(index, num):
    sys.stdout.write("[" + index + " : " + num + "]\r\n")
    sys.stdout.flush()


def write_pass():
    sys.stdout.write('P')
    sys.stdout.flush()


def write_error():
    sys.stdout.write('E')
    sys.stdout.flush()
