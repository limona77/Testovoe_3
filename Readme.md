# 🧊 Тестовое задание


## 🚀 Быстрый старт

1. Склонируй репозиторий

```bash 
  git clone https://github.com/limona77/Testovoe_3
```
2. Собери мигратор
```bash 
  docker build -t migrator .\migrator
 ``` 
3. запусти postgresql
```bash 
  docker-compose up postgres 
``` 
4. запусти миграции

```bash 
docker run  --network host migrator -path=/migrations/  -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/testovoe_3?sslmode=disable {up/down 1}
```
5. запусти программу
```bash
  docker-compose up app
```

## 🌐 GraphQL документация


### Создание поста

- **Запрос**:
  ```
    mutation createPost($input: CreatePostInput!){
    createPost(input: $input){
      id
      title
      content
      allowComments
    }
  }
  ```
- **Variables**:
  ```json
  {
    "input": {
      "title": "Go",
      "content": "Introduction",
      "allowComments": true
    }
  } 
  ```

- **Ответ**:
  ```json
  {
    "data": {
      "createPost": {
        "id": "1",
        "title": "Go",
        "content": "Introduction",
        "allowComments": true
      }
    }
  }
  ```

### Обновление поста

- **Запрос**:
  ```
  mutation updatePost($input: UpdatePostInput!){
    updatePost(input: $input){
      id
      title
      content
      allowComments
    }
  }
  ```
- **Variables**:
  ```json
  {
    "input": {
      "id": 1,
      "title": "React",
      "content": "Introduction",
      "allowComments": true
      }
  } 
  ```

- **Ответ**:
  ```json
  {
  "data": {
      "updatePost": {
        "id": "1",
        "title": "React",
        "content": "Introduction",
        "allowComments": true
      }
    }
  }
  ```

### Создание комментария к посту

- **Запрос**:
  ```
  mutation createComment($input:CreateCommentInput!){
      createComment(input: $input){
      postId
      parentId
      author
      content
    }
  }
  ```
- **Variables**:
  ```json
  {
    "input":{
      "postId": 1,
      "author": "Andrey",
      "content": "✅"
    }
  }
  ```

- **Ответ**:
  ```json
  {
    "data": {
      "createComment": {
        "postId": "1",
        "parentId": null,
        "author": "Andrey",
        "content": "✅"
      }
    }
  }
  ```
### Получение всех постов с комментариями

- **Запрос**:
  ```
  query {
    posts{
      id
      title
      content
      comments{
        postId
        author
        content
      }
    }
  }
  ```
- **Variables**:
  ```json
  {
    "input":{
      "postId": 1,
      "author": "Andrey",
      "content": "✅"
    }
  }
  ```

- **Ответ**:
  ```json
  {
    "data": {
      "posts": [
        {
          "id": "1",
          "title": "React",
          "content": "Introduction",
          "comments": [
            {
              "postId": "1",
              "author": "Andrey",
              "content": "✅"
            }
          ]
        }
      ]
    }
  }
  ```
### Получение всех комментариев

- **Запрос**:
  ```
  query GetComments($postId: ID!, $cursor: Int, $limit: Int){
    comments(postId: $postId, cursor: $cursor, limit: $limit){
      id
      postId
      author
      content
    }
  }

  ```
- **Variables**:
  ```json
  {
    "postId": "2",
    "cursor": 1,
    "limit": 4
  }
  ```

- **Ответ**:
  ```json
  {
    "data": {
      "comments": [
        {
          "id": "2",
          "postId": "1",
          "author": "Andrey",
          "content": "❌"
        },
        {
          "id": "3",
          "postId": "1",
          "author": "Andrey",
          "content": "🧊"
        },
        {
          "id": "4",
          "postId": "1",
          "author": "Andrey",
          "content": "😎"
        },
        {
          "id": "5",
          "postId": "1",
          "author": "Andrey",
          "content": "🚀"
        }
      ]
    }
  }
  ```

### Удаление комментария

- **Запрос**:
  ```
  mutation {
    deleteComment(id: 1)
  }
  ```
- **Ответ**:
  ```json
   {
    "data": {
      "deleteComment": true
    }
   }
  ```

### Удаление поста

- **Запрос**:
  ```
  mutation {
    deletePost(id: 1)
  }
  ```
- **Ответ**:
  ```json
   {
    "data": {
      "deletePost": true
    }
  }
  ```
