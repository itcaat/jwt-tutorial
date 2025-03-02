# Используем официальный Nginx
FROM nginx:1.27.4

# Удаляем стандартную конфигурацию
RUN rm /etc/nginx/conf.d/default.conf

# Копируем наш конфиг
COPY nginx.conf /etc/nginx/conf.d/

# Копируем файлы UI
COPY . /usr/share/nginx/html

# Открываем порт 80
EXPOSE 80