import sys


def write_tail(ctype):
    DATA_GETTING_TIPS = '\r\nFinish Getting data[' + ctype + ']\r\n'
    sys.stdout.write(DATA_GETTING_TIPS)
    sys.stdout.flush()


def write_head(ctype):
    DATA_GETTING_TIPS = 'Getting data[' + ctype + ']:\r\n'
    sys.stdout.write(DATA_GETTING_TIPS)
    sys.stdout.flush()


def write_exec():
    sys.stdout.write('#')
    sys.stdout.flush()


def write_pass():
    sys.stdout.write('*')
    sys.stdout.flush()
