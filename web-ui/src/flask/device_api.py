import os
import requests
import json
import logging

from flask import request, abort
from src.common.objects import SDAManager, Port


class DeviceAPI:
    def __init__(self):
        pass

    @classmethod
    def register_api(cls, app):
        # Get/Set SDA Manager IP
        @app.route("/sdamanager/address", methods=["GET", "POST"])
        def sda_manager_address():
            logging.info("[" + request.method + "] sda manager address - IN")

            if request.method == "GET":
                return SDAManager().get_sda_manager_ip(), 200
            elif request.method == "POST":
                data = json.loads(request.data)
                SDAManager.set_sda_manager_ip(data["ip"])
                return SDAManager().get_sda_manager_ip(), 200
            else:
                return abort(404)

        # Get devices(SDAs) Info.
        @app.route("/sdamanager/devices", methods=["GET"])
        def sda_manager_devices():
            logging.info("[" + request.method + "] sda manager devices - IN")

            l = list()
            ret = dict()
            response = requests.get(
                url="http://" + SDAManager().get_sda_manager_ip() + ":" + str(Port.sda_manager_port()) + "/api/v1/agents",
                timeout=300)

            if response.status_code is not 200:
                logging.error("SDAM Server Return Error, Error Code(" + str(response.status_code) + ") - OUT")
                abort(500)

            for obj in response.json()["agents"]:
                d = dict()
                if "id" in obj:
                    d.update({"id": str(obj["id"])})
                if "host" in obj:
                    d.update({"host": str(obj["host"])})
                if "port" in obj:
                    d.update({"port": str(obj["port"])})
                l.append(d)

            ret.update({"devices": l})
            return json.dumps(json.dumps(ret)), 200

        # Set device(SDA) Info.
        @app.route("/sdamanager/device", methods=["POST"])
        def sda_manager_device():
            logging.info("[" + request.method + "] sda manager device - IN")

            data = json.loads(request.data)
            SDAManager.set_device_id(data["id"])
            SDAManager.set_device_ip(data["host"])
            SDAManager.set_device_port(data["port"])

            return "device", 200

        # Get apps Info.
        @app.route("/sdamanager/apps", methods=["GET"])
        def sda_manager_apps():
            logging.info("[" + request.method + "] sda manager apps - IN")

            l = list()
            apps = list()
            d = dict()
            ret = dict()

            response = requests.get(
                url="http://" + SDAManager().get_sda_manager_ip() + ":" + str(Port.sda_manager_port()) + "/api/v1/agents/"
                    + SDAManager.get_device_id(),
                timeout=300)

            if response.status_code is not 200:
                logging.error("SDAM Server Return Error, Error Code(" + str(response.status_code) + ") - OUT")
                abort(500)

            root_path = os.getcwd()
            with open(root_path + "/static/user/apps", 'r') as content_file:
                content = content_file.read()
                if content != "":
                    apps = json.loads(content)["apps"]

            for obj in response.json()["apps"]:
                d.update({"id": str(obj)})
                response2 = requests.get(
                    url="http://" + SDAManager().get_sda_manager_ip() + ":" + str( Port.sda_manager_port()) + "/api/v1/agents/"
                        + SDAManager.get_device_id() + "/apps/" + str(obj),
                    timeout=300)

                if response2.status_code is not 200:
                    logging.error("SDAM Server Return Error, Error Code(" + str(response.status_code) + ") - OUT")
                    abort(500)
                    
                d.update({"services": len(response2.json()["services"])})
                d.update({"state": response2.json()["state"]})

                for app in apps:
                    if "id" in app and app["id"] == str(obj):
                        d.update({"name": app["name"]})

                l.append(d)

            ret.update({"device": "IP: " + SDAManager.get_device_ip() + ", PORT: " + SDAManager.get_device_port(),
                      "apps": l})

            return json.dumps(json.dumps(ret)), 200

        # Set app Id
        @app.route("/sdamanager/app", methods=["POST", "DELETE"])
        def sda_manager_app():
            logging.info("[" + request.method + "] sda manager app - IN")
            if request.method == "POST":
                data = json.loads(request.data)
                SDAManager.set_app_id(data["id"])
                return "", 200
            elif request.method == "DELETE":
                response = requests.delete(
                    url="http://" + SDAManager().get_sda_manager_ip() + ":" + str(
                        Port.sda_manager_port()) + "/api/v1/agents/"
                        + SDAManager.get_device_id() + "/apps/" + SDAManager.get_app_id(),
                    timeout=300)

                if response.status_code is not 200:
                    logging.error("SDAM Server Return Error, Error Code(" + str(response.status_code) + ") - OUT")
                    abort(500)

                return "", 200
            else:
                return abort(404)

        # Install an app
        @app.route("/sdamanager/app/install", methods=["POST"])
        def sda_manager_app_install():
            logging.info("[" + request.method + "] sda manager app install - IN")

            l = list()
            d = dict()
            data = json.loads(request.data)

            response = requests.post(
                url="http://" + SDAManager().get_sda_manager_ip() + ":" + str(Port.sda_manager_port()) + "/api/v1/agents/"
                    + SDAManager.get_device_id() + "/deploy",
                data=data["data"],
                timeout=300)

            if response.status_code is not 200:
                logging.error("SDAM Server Return Error, Error Code(" + str(response.status_code) + ") - OUT")
                abort(500)

            d.update({"id": response.json()["id"], "name": data["name"]})

            root_path = os.getcwd()
            with open(root_path + "/static/user/apps", 'r') as content_file:
                content = content_file.read()
                if content == "":
                    apps = {"apps": l}
                else:
                    apps = json.loads(content)

            with open(root_path + "/static/user/apps", 'w+') as content_file:
                apps["apps"].append(d)
                content_file.write(json.dumps(apps))

            return "", 200

        # Start an app
        @app.route("/sdamanager/app/start", methods=["GET"])
        def sda_manager_app_start():
            logging.info("[" + request.method + "] sda manager app update - IN")

            response = requests.post(
                url="http://" + SDAManager().get_sda_manager_ip() + ":" + str(Port.sda_manager_port()) + "/api/v1/agents/"
                    + SDAManager.get_device_id() + "/apps/" + SDAManager.get_app_id()
                    + "/start",
                timeout=300)

            if response.status_code is not 200:
                logging.error("SDAM Server Return Error, Error Code(" + str(response.status_code) + ") - OUT")
                abort(500)

            return "", 200

        # Stop an app
        @app.route("/sdamanager/app/stop", methods=["GET"])
        def sda_manager_app_stop():
            logging.info("[" + request.method + "] sda manager app update - IN")

            response = requests.post(
                url="http://" + SDAManager().get_sda_manager_ip() + ":" + str(Port.sda_manager_port()) + "/api/v1/agents/"
                    + SDAManager.get_device_id() + "/apps/" + SDAManager.get_app_id()
                    + "/stop",
                timeout=300)

            if response.status_code is not 200:
                logging.error("SDAM Server Return Error, Error Code(" + str(response.status_code) + ") - OUT")
                abort(500)

            return "", 200

        # Update an app
        @app.route("/sdamanager/app/update", methods=["GET"])
        def sda_manager_app_update():
            logging.info("[" + request.method + "] sda manager app update - IN")

            response = requests.post(
                url="http://" + SDAManager().get_sda_manager_ip() + ":" + str(Port.sda_manager_port()) + "/api/v1/agents/"
                    + SDAManager.get_device_id() + "/apps/" + SDAManager.get_app_id()
                    + "/update",
                timeout=300)

            if response.status_code is not 200:
                logging.error("SDAM Server Return Error, Error Code(" + str(response.status_code) + ") - OUT")
                abort(500)

            return "", 200

        # Get/Update app Yaml file to SDA DB
        @app.route("/sdamanager/app/yaml", methods=["GET", "POST"])
        def sda_manager_app_yaml():
            logging.info("[" + request.method + "] sda manager app YAML - IN")

            if request.method == "GET":
                d = dict()
                response = requests.get(
                    url="http://" + SDAManager().get_sda_manager_ip() + ":" + str(Port.sda_manager_port()) + "/api/v1/agents/"
                        + SDAManager.get_device_id() + "/apps/" + SDAManager.get_app_id(),
                    timeout=300)

                if response.status_code is not 200:
                    logging.error("SDAM Server Return Error, Error Code(" + str(response.status_code) + ") - OUT")
                    abort(500)

                if "description" in response.json():
                    d = response.json()["description"]

                return json.dumps(json.dumps(d)), 200
            elif request.method == "POST":
                data = json.loads(request.data)

                response = requests.post(
                    url="http://" + SDAManager().get_sda_manager_ip() + ":" + str(Port.sda_manager_port()) + "/api/v1/agents/"
                        + SDAManager.get_device_id() + "/apps/" + SDAManager.get_app_id(),
                    data=data["data"],
                    timeout=300)

                if response.status_code is not 200:
                    logging.error("SDAM Server Return Error, Error Code(" + str(response.status_code) + ") - OUT")
                    abort(500)

                return "", 200

        # Get/Write Yaml file to Web client DB
        @app.route("/sdamanager/yaml", methods=["GET", "POST"])
        def sda_manager_yaml():
            logging.info("[" + request.method + "] sda manager YAML - IN")

            if request.method == "GET":
                root_path = os.getcwd()
                with open(root_path + "/static/user/yamls", 'r') as content_file:
                    content = content_file.read()

                return json.dumps(content), 200
            elif request.method == "POST":
                root_path = os.getcwd()
                with open(root_path + "/static/user/yamls", 'r') as content_file:
                    content = content_file.read()
                    yamls = json.loads(content)
                    yamls["yamls"].append(json.loads(request.data))

                with open(root_path + "/static/user/yamls", 'w+') as content_file:
                    content_file.write(json.dumps(yamls))

                return "", 200
            else:
                logging.error("Unknown Method - OUT")
                return abort(404)

        # Register Device
        @app.route("/sdamanager/register", methods=["POST"])
        def sda_manager_device_register():
            logging.info("[" + request.method + "] sda manager device register - IN")

            d = dict()
            d2 = dict()
            data = json.loads(request.data)

            d2.update({"interval": str(data["interval"])})
            d.update({"ip": str(SDAManager().get_sda_manager_ip() + ":" + str(Port.sda_manager_port())).split(':')[0], "healthCheck": d2})

            response = requests.post(
                url="http://" + data["ip"] + ":" + str(Port.sda_port()) + "/api/v1/register",
                data=json.dumps(d),
                timeout=300)

            if response.status_code is not 200:
                logging.error("SDAM Server Return Error, Error Code(" + str(response.status_code) + ") - OUT")
                abort(500)

            return "", 200

        # Unregister Device
        @app.route("/sdamanager/unregister", methods=["POST"])
        def sda_manager_device_unregister():
            logging.info("[" + request.method + "] sda manager device register - IN")

            response = requests.post(
                url="http://" + SDAManager.get_device_ip() + ":" + SDAManager.get_device_port() + "/api/v1/unregister",
                timeout=300)

            if response.status_code is not 200:
                logging.error("SDAM Server Return Error, Error Code(" + str(response.status_code) + ") - OUT")
                abort(500)

            return "", 200
