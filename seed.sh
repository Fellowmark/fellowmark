#!/bin/sh

ADMIN_TOKEN=$(curl -X GET "http://localhost:5000/admin/auth/login?email=admin@local.com&password=admin" | jq -r '.message')

PROF_ID=$(curl -X POST -H "Content-Type: application/json" http://localhost:5000/staff/auth/signup -d '{"email": "akshay@fm.com", "password": "12345678", "name": "Akshay Narayan"}' | jq '.ID')

MODULE_ID=$(curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:5000/module -d '{"code": "CS2103", "name":"Software Engineering", "semester": "AY2021/2022 1"}' | jq '.ID')
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:5000/module/supervise -d '{"moduleId": '"$MODULE_ID"', "staffId":'"$PROF_ID"'}'

STUDENT_ID=$(curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "gerard@fm.com", "password": "12345678", "name": "Gerard Berg"}' | jq '.ID')
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:5000/module/enroll -d '{"moduleId": '"$MODULE_ID"', "studentId":'"$STUDENT_ID"'}'

STUDENT_ID=$(curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "linwood@fm.com", "password": "12345678", "name": "Linwood Clark"}' | jq '.ID')
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:5000/module/enroll -d '{"moduleId": '"$MODULE_ID"', "studentId":'"$STUDENT_ID"'}'

STUDENT_ID=$(curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "derek@fm.com", "password": "12345678", "name": "Derek Marquez"}' | jq '.ID')
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:5000/module/enroll -d '{"moduleId": '"$MODULE_ID"', "studentId":'"$STUDENT_ID"'}'

STUDENT_ID=$(curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "romeo@fm.com", "password": "12345678", "name": "Romeo Estrada"}' | jq '.ID')
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:5000/module/enroll -d '{"moduleId": '"$MODULE_ID"', "studentId":'"$STUDENT_ID"'}'

STUDENT_ID=$(curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "alfredo@fm.com", "password": "12345678", "name": "Alfredo Schultz"}' | jq '.ID')
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:5000/module/enroll -d '{"moduleId": '"$MODULE_ID"', "studentId":'"$STUDENT_ID"'}'

STUDENT_ID=$(curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "lesley@fm.com", "password": "12345678", "name": "Lesley Stewart"}' | jq '.ID')
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:5000/module/enroll -d '{"moduleId": '"$MODULE_ID"', "studentId":'"$STUDENT_ID"'}'

STUDENT_ID=$(curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "alton@fm.com", "password": "12345678", "name": "Alton Shea"}' | jq '.ID')
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:5000/module/enroll -d '{"moduleId": '"$MODULE_ID"', "studentId":'"$STUDENT_ID"'}'

STUDENT_ID=$(curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "suzanne@fm.com", "password": "12345678", "name": "Suzanne Smith"}' | jq '.ID')
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:5000/module/enroll -d '{"moduleId": '"$MODULE_ID"', "studentId":'"$STUDENT_ID"'}'

STUDENT_ID=$(curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "ervin@fm.com", "password": "12345678", "name": "Ervin Bass"}' | jq '.ID')
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:5000/module/enroll -d '{"moduleId": '"$MODULE_ID"', "studentId":'"$STUDENT_ID"'}'

STUDENT_ID=$(curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "vance@fm.com", "password": "12345678", "name": "Vance Kent"}' | jq '.ID')
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:5000/module/enroll -d '{"moduleId": '"$MODULE_ID"', "studentId":'"$STUDENT_ID"'}'

STUDENT_ID=$(curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "frances@fm.com", "password": "12345678", "name": "Frances Erickson"}' | jq '.ID')
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:5000/module/enroll -d '{"moduleId": '"$MODULE_ID"', "studentId":'"$STUDENT_ID"'}'

STUDENT_ID=$(curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "trisha@fm.com", "password": "12345678", "name": "Trisha Lin"}' | jq '.ID')
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:5000/module/enroll -d '{"moduleId": '"$MODULE_ID"', "studentId":'"$STUDENT_ID"'}'

STUDENT_ID=$(curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "ana@fm.com", "password": "12345678", "name": "Ana Franklin"}' | jq '.ID')
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:5000/module/enroll -d '{"moduleId": '"$MODULE_ID"', "studentId":'"$STUDENT_ID"'}'

STUDENT_ID=$(curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "christi@fm.com", "password": "12345678", "name": "Christi Hubbard"}' | jq '.ID')
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:5000/module/enroll -d '{"moduleId": '"$MODULE_ID"', "studentId":'"$STUDENT_ID"'}'

STUDENT_ID=$(curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "herschel@fm.com", "password": "12345678", "name": "Herschel Snyder"}' | jq '.ID')
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:5000/module/enroll -d '{"moduleId": '"$MODULE_ID"', "studentId":'"$STUDENT_ID"'}'

STUDENT_ID=$(curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "dallas@fm.com", "password": "12345678", "name": "Dallas Brennan"}' | jq '.ID')
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:5000/module/enroll -d '{"moduleId": '"$MODULE_ID"', "studentId":'"$STUDENT_ID"'}'

STUDENT_ID=$(curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "yesenia@fm.com", "password": "12345678", "name": "Yesenia Schroeder"}' | jq '.ID')
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:5000/module/enroll -d '{"moduleId": '"$MODULE_ID"', "studentId":'"$STUDENT_ID"'}'

STUDENT_ID=$(curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "mitch@fm.com", "password": "12345678", "name": "Mitch Larsen"}' | jq '.ID')
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:5000/module/enroll -d '{"moduleId": '"$MODULE_ID"', "studentId":'"$STUDENT_ID"'}'

STUDENT_ID=$(curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "margarito@fm.com", "password": "12345678", "name": "Margarito Munoz"}' | jq '.ID')
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:5000/module/enroll -d '{"moduleId": '"$MODULE_ID"', "studentId":'"$STUDENT_ID"'}'

STUDENT_ID=$(curl -X POST -H "Content-Type: application/json" http://localhost:5000/student/auth/signup -d '{"email": "violet@fm.com", "password": "12345678", "name": "Violet Liu"}' | jq '.ID')
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:5000/module/enroll -d '{"moduleId": '"$MODULE_ID"', "studentId":'"$STUDENT_ID"'}'

ASSIGNMENT_ID=$(curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:5000/assignment -d '{"name": "Assignment 1", "moduleId": '"$MODULE_ID"', "groupSize": 3}' | jq '.ID')

curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:5000/assignment/pairs/initialize -d '{"id": '"$ASSIGNMENT_ID"'}'
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:5000/assignment/pairs/assign -d '{"id": '"$ASSIGNMENT_ID"'}'

QUESTION_ID=$(curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:5000/assignment/question -d '{"questionNumber": 1, "assignmentId": '"$ASSIGNMENT_ID"', "questionText": "This is a question"}' | jq '.ID')

curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:5000/assignment/rubric -d '{"questionId": '"$QUESTION_ID"', "criteria": "Reliability", "description": "How reliable is the system", "minMark": 1, "maxMark": 10}'
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:5000/assignment/rubric -d '{"questionId": '"$QUESTION_ID"', "criteria": "Quality", "description": "What is the quality of your system", "minMark": 1, "maxMark": 10}'
