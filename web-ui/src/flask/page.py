from flask import render_template, request, abort


class Page:
    def __init__(self):
        pass

    @classmethod
    def register_page(cls, app):
        # Pages
        @app.route("/", methods=["GET"])
        def index():
            if request.method == "GET":
                return render_template("index.html")
            else:
                return abort(404)

        @app.route("/sdamanager", methods=["GET"])
        def settings():
            if request.method == "GET":
                return render_template("sda_manager.html")
            else:
                return abort(404)

        @app.route("/contact", methods=["GET"])
        def contact():
            if request.method == "GET":
                return render_template("contact.html")
            else:
                return abort(404)

        @app.route("/device", methods=["GET"])
        def device():
            if request.method == "GET":
                return render_template("device.html")
            else:
                return abort(404)

        @app.route("/app", methods=["GET"])
        def app():
            if request.method == "GET":
                return render_template("app.html")
            else:
                return abort(404)
