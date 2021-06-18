Use this command to add a mosque to a database with some common data
Replace {mosqueid} with the mosqueid you want to create a DB for
mongoimport --jsonArray --db {mosqueid} --collection "01-01-2021"  --file ~/Entrycontrol-API/repos/creation.json
