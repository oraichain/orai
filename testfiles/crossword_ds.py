#!/usr/bin/python3
import requests

def main():
    r = requests.get('http://128.199.241.140:8085/v1/gen_crossword').json()
    status = r['data']['status']
    data = r['data']
    if status == 200:
        response = {'question_matrix': data['question_matrix'], 'answer_matrix': data['answer_matrix'], 'questions_list': data['questions_list']}
        return response
    return None

if __name__ == "__main__":
    print(main())