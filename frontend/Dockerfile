# Use the official Nginx image
FROM nginx:latest

# Set the working directory inside the container
WORKDIR /usr/share/nginx/html

# Copy all frontend files to the Nginx HTML directory
COPY ./ /usr/share/nginx/html

# Expose port 80 for the Nginx server
EXPOSE 80

# Command to run Nginx
CMD ["nginx", "-g", "daemon off;"]
