import os


def import_env():
    if os.path.exists('.env'):
        for line in open('.env'):
            var = line.strip().split('=')
            if len(var) == 2:
                key, value = var[0].strip(), var[1].strip()
                os.environ[key] = value
