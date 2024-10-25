import json
import os

print("OPEN NEW GENESIS")
genesis = open(os.environ["NEW_GENESIS"], "r").read()
genesis_json_object = json.loads(genesis)

print("OPEN OLD GENESIS")
exported_genesis = open(os.environ["OLD_GENESIS"], "r").read()
exported_genesis_json_object = json.loads(exported_genesis)

print("PULL STATE OUT OF OLD GENESIS")
crosschain = exported_genesis_json_object["app_state"]["crosschain"]
observer = exported_genesis_json_object["app_state"]["observer"]
emissions = exported_genesis_json_object["app_state"]["emissions"]
fungible = exported_genesis_json_object["app_state"]["fungible"]
evm = exported_genesis_json_object["app_state"]["evm"]
auth_accounts = exported_genesis_json_object["app_state"]["auth"]["accounts"]

print("MANIPULATE NEW GENESIS")
genesis_json_object["app_state"]["auth"]["accounts"] = genesis_json_object["app_state"]["auth"]["accounts"] + auth_accounts
genesis_json_object["app_state"]["crosschain"] = crosschain
genesis_json_object["app_state"]["observer"] = observer
genesis_json_object["app_state"]["emissions"] = emissions
genesis_json_object["app_state"]["fungible"] = fungible

evm_accounts = []
for index, account in enumerate(evm["accounts"]):
    if account["address"] == "0x0000000000000000000000000000000000000001":
        print("pop account", account["address"])
    elif account["address"] == "0x0000000000000000000000000000000000000006":
        print("pop account", account["address"])
    elif account["address"] == "0x0000000000000000000000000000000000000002":
        print("pop account", account["address"])
    elif account["address"] == "0x0000000000000000000000000000000000000002":
        print("pop account", account["address"])
    elif account["address"] == "0x0000000000000000000000000000000000000008":
        print("pop account", account["address"])
    else:
        evm_accounts.append(account)
evm["accounts"] = evm_accounts
genesis_json_object["app_state"]["evm"] = evm

print("WRITE GENESIS-EDITED")
genesis = open("genesis-edited.json", "w")
genesis_string = json.dumps(genesis_json_object, indent=2)
dumped_genesis_object = genesis_string.replace("0x0000000000000000000000000000000000000001","0x387A12B28fe02DcAa467c6a1070D19B82F718Bb5")
genesis.write(genesis_string)
genesis.close()
