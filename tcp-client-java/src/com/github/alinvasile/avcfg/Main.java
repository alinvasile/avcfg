package com.github.alinvasile.avcfg;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.Socket;
import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.util.concurrent.TimeUnit;

public class Main {

    static final byte API_VERSION = 0x01 ;

    public static void main(String[] args) throws IOException {
        String host = args[0];
        String port = args[1];
        String property = args[2];

        System.out.println("API version: " + API_VERSION);

        long start = System.nanoTime();
        Socket socket = new Socket(host, Integer.parseInt(port));
        System.out.println("Connection time: " + TimeUnit.MILLISECONDS.convert(System.nanoTime() - start, TimeUnit.NANOSECONDS) + " ms");

        start = System.nanoTime();
        OutputStream socketOutputStream = socket.getOutputStream();

        ByteBuffer b = ByteBuffer.allocate(2);
        b.order(ByteOrder.BIG_ENDIAN);
        b.putShort(API_VERSION);
        socketOutputStream.write(b.array());

        byte[] bytes = property.getBytes();
        b = ByteBuffer.allocate(4);
        b.order(ByteOrder.BIG_ENDIAN);
        b.putInt(bytes.length );
        socketOutputStream.write(b.array());

//        System.out.println("Size of string " + bytes.length);

        socketOutputStream.write(bytes);

        System.out.println("Write time: " + TimeUnit.MILLISECONDS.convert(System.nanoTime() - start, TimeUnit.NANOSECONDS) + " ms");

        readResponse(socket);



    }

    private static void readResponse(Socket socket) throws IOException {
        long start;
        start = System.nanoTime();

        InputStream inputStream = socket.getInputStream();

        byte[] version = new byte[2];
        inputStream.read(version);

        ByteBuffer wrapped = ByteBuffer.wrap(version);
//        System.out.println("Api version: " + wrapped.getShort());

        byte[] length = new byte[4];
        inputStream.read(length);
        ByteBuffer w2 = ByteBuffer.wrap(length);

        int lng = w2.getInt();
//        System.out.println("Length: " + lng);

        byte[] response = new byte[lng * 2];
        inputStream.read(response);
        System.out.println("Response time: " + TimeUnit.MILLISECONDS.convert(System.nanoTime() - start, TimeUnit.NANOSECONDS) + " ms");

        System.out.println("Response: " + new String(response));
    }
}
