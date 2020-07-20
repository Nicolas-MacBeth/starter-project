package foodcounter;

import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import com.sun.net.httpserver.HttpServer;

public class Test {

    private static int count = 0;
    private static int port = 8084;

    // Simple http server in java to showcase JVM metrics, it simply counts the amount of requests to a certain endpoint
    public static void main(String[] args) throws Exception {
        HttpServer server = HttpServer.create(new InetSocketAddress(port), 0);
        server.createContext("/count", new CounterHandler());
        server.setExecutor(null); // creates a default executor
        System.out.println("[Counter] Starting Java server to count requests, listening on port: " + port);
        server.start();
    }

    static class CounterHandler implements HttpHandler {
        @Override
        public void handle(HttpExchange t) throws IOException {
            String response = "Hello World!";
            t.sendResponseHeaders(200, response.length());
            OutputStream os = t.getResponseBody();
            os.write(response.getBytes());
            os.close();
            System.out.println("Request count from Java server: " + ++count);
        }
    }
}