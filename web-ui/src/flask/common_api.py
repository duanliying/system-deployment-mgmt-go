import logging

from flask import request
from src.common.objects import SDAManager


class CommonAPI:
    def __init__(self):
        pass

    @classmethod
    def register_api(cls, app):
        @app.route("/sdamanager/check/address", methods=["GET"])
        def sda_manager_check_address():
            logging.info("[" + request.method + "] sda manager check address - IN")

            if SDAManager.get_sda_manager_address() == "":
                logging.error("No SDAM Address - OUT")
                return "", 500
            else:
                return SDAManager.get_sda_manager_address(), 200
