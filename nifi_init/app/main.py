import json
import nipyapi
from dataclasses import dataclass
import os

@dataclass
class Settings():
    nifi_host: str
    nifi_registry_host: str
    spec_path: str

    def prepare_settings():
        return Settings(
            os.getenv("NIFI_HOST", "http://nifi:8080/nifi-api"),
            os.getenv("NIFI_REGISTRY_HOST", "http://nifi-registry:18080/nifi-registry-api"),
            os.getenv("SPEC_PATH", "/spec/init_pg_spec.json")
        )

def is_registry_connected(registry_name):
    l = nipyapi.versioning.list_registry_clients()
    if not l:
        return False
    for r in l.registries:
        if r.component.name == registry_name:
            return True
    return False

def get_registry_identifier(registry_name):
    l = nipyapi.versioning.list_registry_clients()
    if not l:
        return None
    for r in l.registries:
        if r.component.name == registry_name:
            return r.component.id
    return None

def create_registry_if_not_exist(name, uri, description="auto generated registry"):
    reg_id = get_registry_identifier(name)
    if reg_id is not None:
        print(f"registry {name} already exist with id {reg_id}")
        return reg_id
    reg = nipyapi.versioning.create_registry_client(name, uri, description)
    return reg.component.id

def is_pg_exist(name):
    l = nipyapi.canvas.list_all_process_groups()
    for pg in l:
        if pg.component.name == name:
            return True
    return False

def get_pg(name):
    l = nipyapi.canvas.list_all_process_groups()
    for pg in l:
        if pg.component.name == name:
            return pg
    return None

def get_bucket_identifier(bucket_name):
    buckets = nipyapi.versioning.list_registry_buckets()
    for b in buckets:
        if b.name == bucket_name:
            return b.identifier
    return None

def get_flow_identifier(flow_name, bucket_id):
    flows = nipyapi.versioning.list_flows_in_bucket(bucket_id)
    for f in flows:
        if f.name == flow_name:
            return f.identifier
    return None

def get_processor(name, pg_id):
    processors = nipyapi.canvas.list_all_processors(pg_id)
    for proc in processors:
        if proc.component.name == name:
            return proc
    return None

if __name__ == '__main__':
    s = Settings.prepare_settings()

    nipyapi.config.nifi_config.host = s.nifi_host
    nipyapi.config.registry_config.host = s.nifi_registry_host

    nifi_check_endpoint = '-'.join(s.nifi_host.split('-')[:-1])
    print(f"nifi check endpoint: '{nifi_check_endpoint}'")
    nifi_reg_check_endpoint = '-'.join(s.nifi_registry_host.split('-')[:-1])
    print(f"nifi registry check endpoint: '{nifi_reg_check_endpoint}'")

    print("checking nifi")
    nipyapi.utils.wait_to_complete(
        test_function=nipyapi.utils.is_endpoint_up,
        endpoint_url=nifi_check_endpoint,
        nipyapi_delay=nipyapi.config.long_retry_delay,
        nipyapi_max_wait=nipyapi.config.long_max_wait
    )

    print("checking registry")
    nipyapi.utils.wait_to_complete(
        test_function=nipyapi.utils.is_endpoint_up,
        endpoint_url=nifi_reg_check_endpoint,
        nipyapi_delay=nipyapi.config.long_retry_delay,
        nipyapi_max_wait=nipyapi.config.long_max_wait
    )

    # nipyapi.security.service_login(service='nifi', username='nifi', password='nifi123456789')

    x_step = 20
    x = 0

    conf = json.load(open('/spec/init_pg_spec.json'))
    for registry in conf:
        reg_id = create_registry_if_not_exist(
            registry['registry_name'],
            registry['url'],
            "auto generated registry"
        )
        for bucket in registry['buckets']:
            buck_id = get_bucket_identifier(bucket['name'])
            if buck_id is None:
                print(f"cant find bucket {bucket['name']} in registry {registry['registry_name']}")
                continue
            for pg in bucket['process_groups']:
                flow_id = get_flow_identifier(pg['pg_name'], buck_id)
                if flow_id is None:
                    print(f"cant find flow {pg['pg_name']} in bucket {bucket['name']} in registry {registry['registry_name']}")
                    continue
                _pg = get_pg(pg['pg_name'])
                if _pg is None:
                    _pg = nipyapi.versioning.deploy_flow_version(
                        nipyapi.canvas.get_root_pg_id(),
                        (x,0),
                        buck_id,
                        flow_id,
                        reg_id
                    )
                    x += x_step

                nipyapi.canvas.schedule_process_group(_pg.id, scheduled=False)
                nipyapi.versioning.update_flow_ver(_pg, target_version=None)
                for setup in pg['setups']:
                    processor = get_processor(setup['processor_name'], _pg.id)
                    if processor is None:
                        print(f"cant find processor {setup['name']} in processor group {pg['pg_name']} in bucket {bucket['name']} in registry {registry['registry_name']}")
                        continue
                    if setup.get('properties', []):
                        properties = processor.component.config.properties
                        for prop in setup['properties']:
                            properties[prop['name']] = prop['value']
                        config = processor.component.config
                        config.properties = properties
                        nipyapi.canvas.update_processor(processor, config)
                for controller in nipyapi.canvas.list_all_controllers(pg_id=_pg.id, descendants=True):
                    nipyapi.canvas.schedule_controller(controller, scheduled=True, refresh=True)
                nipyapi.canvas.schedule_process_group(_pg.id, scheduled=True)