# school
housing=# GRANT SELECT, UPDATE ON users TO auth_user;
<br>
housing=# GRANT SELECT, UPDATE ON users_id_seq TO auth_user;
<br>
housing=# GRANT SELECT, UPDATE, INSERT ON passports_info TO auth_user;
<br>
housing=# GRANT SELECT, UPDATE ON passports_info_id_seq TO auth_user;
<br>
protoc -I pkg/authService/ --go_out=plugins=grpc:pkg/authService/ pkg/authService/authorize.proto