import os
from flask import Flask
from flask_sqlalchemy import SQLAlchemy
from sqlalchemy import Column, Integer, BigInteger, String

class Config:
  DB_USER = os.environ['DB_USER']
  DB_PASSWORD = os.environ['DB_PASSWORD']
  DB_NAME = os.environ['DB_NAME']
  DB_HOST = os.environ['DB_HOST']
  DB_PORT = os.environ['DB_PORT']

  SQLALCHEMY_DATABASE_URI = f'postgresql+psycopg2://{DB_USER}:{DB_PASSWORD}@{DB_HOST}:{DB_PORT}/{DB_NAME}'
  SQLALCHEMY_POOL_SIZE = 100
  SQLALCHEMY_TRACK_MODIFICATIONS = False

app = Flask(__name__)
app.config.from_object(Config)
db = SQLAlchemy(app)

class User(db.Model):
  __tablename__ = 'users'

  id = Column(BigInteger, primary_key=True)
  name = Column(String)
  age = Column(Integer)

@app.route('/', methods=['GET'])
def index():
  user = User(name='Test!', age=1)
  db.session.add(user)
  db.session.commit()

  user_db = User.query.get(user.id)
  user_db.age += 1
  db.session.commit()

  db.session.delete(user_db)
  db.session.commit()

  return {'id': user_db.id, 'name': user_db.name, 'age': user_db.age}
