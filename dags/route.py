from scripts import bash, test
import os
SCRIPT = os.environ.get("script")


if SCRIPT == "test":
    test.start()
elif SCRIPT == "bash":
    bash.start()
