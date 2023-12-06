# edisco-test
test and validate eDiscovery tools

## Unpack `pst` files

I am using `readpst` in a Debian container image I built for the `batch-decipher-pst` project.

This command will unpack the pst file with the following options.

- a = only save *attachments* that match this list of extensions, `.fubar` is meant to not match **anything** in effect don't save attachments
- D = include *deleted* items
- b = don't save RTF *body*
- S = *separate* files format means save attachments and RTF body separately from the email
- t = *type* of messages to unpack `e` = email
- j = number of threads to run, `nproc` command prints the number of CPUs you have
- o = *output* directory path, output will be a folder tree with numbered files, no extensions

The last argument is the path for the input pst file.

```bash
podman run -it --rm -v $(pwd):/input:z batch-decipher-pst_busybee:latest readpst -a ".fubar" -DbS -t e -j $(nproc) -o /input/allIn /input/in1.pst
```

## Ingest `eml` files

Run the ingest command to parse email files into a newline-delimited json file.

```bash
edisco-test ingestEmail --in-dir output-path-from-readpst --out all1.json
```

## Postgres

Run the pg container.

```bash
podman run --name postgres -u postgres -d -v $(pwd):/input:z -v postgres:/var/lib/posgresql/data -e POSTGRES_PASSWORD="fuar" postgres:latest
```

Connect to the db using `psql`.

```bash
podman exec -it -u postgres postgres psql
```

Ingest the json data into a Postgres Database for analysis.

```sql
-- ingest json file
create table temp(data jsonb); 
\copy temp (data) from program 'sed -e ''s/\\/\\\\/g'' /input/all1.json'

-- create table to receive json data
create table eml_all (eml_from text, eml_to text, subject text, eml_date timestamptz not null);

-- load json data into table
insert into eml_all select data->>'From' as eml_from, data->>'To' as eml_to, data->>'Subject' as subject, (data->>'Date')::timestamptz as eml_date from temp;

-- cleanup
drop table temp;
```

## Reference Query

Build a sql query that matches the logic of the conditions you are applying in the system under test.
Apply this query to the reference data set and save in a variable.

```sql
select * into reference from eml_all where eml_from = 'jane.doe@local';
```

## Ingest deliverable from the system under test

Follow the same procedures as ingesting the reference data set.
Save it to a table named `test`.
Query all data from the test table.

```sql
select * into test from test_from_conditions;
```

## Find the difference between the reference and test

```sql
(select * from reference) EXCEPT (select * from test);
```

This shows all records that exist in `reference` that don't exist in `test`.
Now do the other way arround.

```sql
(select * from test) EXCEPT (select * from reference);
```