import os
from fastapi import FastAPI
from sqlalchemy import create_engine, Column, String, BigInteger, Integer
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker

class Config:
  DB_USER = os.environ['DB_USER']
  DB_PASSWORD = os.environ['DB_PASSWORD']
  DB_NAME = os.environ['DB_NAME']
  DB_HOST = os.environ['DB_HOST']
  DB_PORT = os.environ['DB_PORT']

  DATABASE_URI = f'postgresql+psycopg2://{DB_USER}:{DB_PASSWORD}@{DB_HOST}:{DB_PORT}/{DB_NAME}'



engine = create_engine(Config.DATABASE_URI, pool_size=100)

SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

Base = declarative_base()

class User(Base):
  __tablename__ = 'users'

  id = Column(BigInteger, primary_key=True)
  name = Column(String)
  age = Column(Integer)


app = FastAPI()

@app.get("/")
def index():
  with SessionLocal() as session:
    user = User(name='Test!', age=1)
    session.add(user)
    session.commit()

    user_db = session.query(User).get(user.id)
    user_db.age += 1
    session.commit()

    session.delete(user_db)
    session.commit()

  return {'id': user_db.id, 'name': user_db.name, 'age': user_db.age}

