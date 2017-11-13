import logging
import argparse

from flask import Flask

from src.common.objects import Port
from src.flask.device_api import DeviceAPI
from src.flask.group_api import GroupAPI
from src.flask.common_api import CommonAPI
from src.flask.page import Page

logging.basicConfig(format='%(levelname)s[%(asctime)s] SDA Web Client: %(message)s', level=logging.DEBUG)

app = Flask(__name__)


def main():
    # run server with bellow arguments
    # Refer this site= https://docs.python.org/2/howto/argparse.html
    parser = argparse.ArgumentParser(
        description="This is Web Client APP for Service Deployment Agent & Service Deployment Agent Manager.\n"
                    "If you have problem, please contact us. (jihun.ha@samsung.com)",
        formatter_class=argparse.RawTextHelpFormatter
    )

    parser.add_argument(
        "-address",
        type=str, default="0.0.0.0:" + Port.client_port(),
        help="The ip address & port you want to open.\n"
             "If you write \"0.0.0.0:<port>\", you can make the server publicly available\n"
             "(default: \"0.0.0.0:"+Port.client_port() + "\")\n"
             " "
    )

    parser.add_argument(
        "-debug",
        type=int, default=0,
        help="debug mode: 1, release mode: 0\n"
        "(default: \"0\")\n"
        " "
    )

    args = parser.parse_args()
    address = args.address.split(":")
    host = address[0]
    port = int(address[1])
    debug = args.debug

    runserver(host, port, debug)


def runserver(host="0.0.0.0", port=5555, debug=True):
    # Rendered Page define
    Page.register_page(app)

    # API define using Web Client
    DeviceAPI.register_api(app)
    GroupAPI.register_api(app)
    CommonAPI.register_api(app)

    # run application with parameters
    app.run(host, port, debug)


if __name__ == "__main__":
    main()
