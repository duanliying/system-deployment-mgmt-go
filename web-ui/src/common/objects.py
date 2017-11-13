### Define Exposed Posts of Web Client ###
class Port():
    @staticmethod
    def sda_manager_port():
        return "48099"

    @staticmethod
    def sda_port():
        return "48098"

    @staticmethod
    def client_port():
        return "5555"


### Define Global Values ###
class SDAManager:
    sda_manager_ip = ""
    group_id = ""
    current_group_name = ""
    device_id = ""
    device_ip = ""
    device_port = ""
    app_id = ""

    def __init__(self):
        pass

    @classmethod
    def set_sda_manager_ip(cls, address):
        cls.sda_manager_ip = address

    @classmethod
    def get_sda_manager_ip(cls):
        return cls.sda_manager_ip

    @classmethod
    def set_group_id(cls, id):
        cls.group_id = id

    @classmethod
    def get_group_id(cls):
        return cls.group_id

    @classmethod
    def set_current_group_name(cls, name):
        cls.current_group_name = name

    @classmethod
    def get_current_group_name(cls):
        return cls.current_group_name

    @classmethod
    def set_device_id(cls, id):
        cls.device_id = id

    @classmethod
    def get_device_id(cls):
        return cls.device_id

    @classmethod
    def set_device_ip(cls, ip):
        cls.device_ip = ip

    @classmethod
    def get_device_ip(cls):
        return cls.device_ip

    @classmethod
    def set_device_port(cls, port):
        cls.device_port = port

    @classmethod
    def get_device_port(cls):
        return cls.device_port

    @classmethod
    def set_app_id(cls, id):
        cls.app_id = id

    @classmethod
    def get_app_id(cls):
        return cls.app_id
