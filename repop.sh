#!/bin/sh

# note professor id
curl -X POST -H "Content-Type: application/json" http://localhost:5000/admin/staff -d '{"email": "akshay@fm.com", "password": "12345678", "name": "Akshay Narayan"}'

curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "gerard@fm.com", "password": "12345678", "name": "Gerard Berg"}'

# note the first id created for gerard.
curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "linwood@fm.com", "password": "12345678", "name": "Linwood Clark"}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "derek@fm.com", "password": "12345678", "name": "Derek Marquez"}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "romeo@fm.com", "password": "12345678", "name": "Romeo Estrada"}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "alfredo@fm.com", "password": "12345678", "name": "Alfredo Schultz"}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "lesley@fm.com", "password": "12345678", "name": "Lesley Stewart"}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "alton@fm.com", "password": "12345678", "name": "Alton Shea"}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "suzanne@fm.com", "password": "12345678", "name": "Suzanne Smith"}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "ervin@fm.com", "password": "12345678", "name": "Ervin Bass"}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "vance@fm.com", "password": "12345678", "name": "Vance Kent"}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "frances@fm.com", "password": "12345678", "name": "Frances Erickson"}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "trisha@fm.com", "password": "12345678", "name": "Trisha Lin"}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "ana@fm.com", "password": "12345678", "name": "Ana Franklin"}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "christi@fm.com", "password": "12345678", "name": "Christi Hubbard"}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "herschel@fm.com", "password": "12345678", "name": "Herschel Snyder"}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "dallas@fm.com", "password": "12345678", "name": "Dallas Brennan"}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "yesenia@fm.com", "password": "12345678", "name": "Yesenia Schroeder"}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "mitch@fm.com", "password": "12345678", "name": "Mitch Larsen"}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "margarito@fm.com", "password": "12345678", "name": "Margarito Munoz"}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "violet@fm.com", "password": "12345678", "name": "Violet Liu"}'


# look at module id created. 
curl -X POST -H "Content-Type: application/json" http://localhost:5000/module -d '{"code": "CS2103", "name":"Software Engineering", "semester": "AY2021/2022 1"}'

# staffId should add id returned on staff creation
curl -X POST -H "Content-Type: application/json" http://localhost:5000/module/supervise -d '{"moduleId": 2, "staffId":4}'

# studentId should match the ids returned by signups
curl -X POST -H "Content-Type: application/json" http://localhost:5000/module/enroll -d '{"moduleId": 2, "studentId":8}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/module/enroll -d '{"moduleId": 2, "studentId":9}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/module/enroll -d '{"moduleId": 2, "studentId":10}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/module/enroll -d '{"moduleId": 2, "studentId":11}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/module/enroll -d '{"moduleId": 2, "studentId":12}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/module/enroll -d '{"moduleId": 2, "studentId":13}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/module/enroll -d '{"moduleId": 2, "studentId":14}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/module/enroll -d '{"moduleId": 2, "studentId":15}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/module/enroll -d '{"moduleId": 2, "studentId":16}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/module/enroll -d '{"moduleId": 2, "studentId":17}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/module/enroll -d '{"moduleId": 2, "studentId":18}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/module/enroll -d '{"moduleId": 2, "studentId":19}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/module/enroll -d '{"moduleId": 2, "studentId":20}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/module/enroll -d '{"moduleId": 2, "studentId":21}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/module/enroll -d '{"moduleId": 2, "studentId":22}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/module/enroll -d '{"moduleId": 2, "studentId":23}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/module/enroll -d '{"moduleId": 2, "studentId":24}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/module/enroll -d '{"moduleId": 2, "studentId":25}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/module/enroll -d '{"moduleId": 2, "studentId":26}'
curl -X POST -H "Content-Type: application/json" http://localhost:5000/module/enroll -d '{"moduleId": 2, "studentId":27}'
