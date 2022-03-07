import java.io.*;

public class JarMain {

    private static String executeCommand(String command) {

        StringBuilder output = new StringBuilder();

        try {
            Process p = Runtime.getRuntime().exec(command);
            p.waitFor();
            BufferedReader reader = new BufferedReader(new InputStreamReader(p.getInputStream()));

            String line;
            while((line = reader.readLine()) != null) {
                output.append(line).append("\n");
            }
        } catch (Exception e) {
            return e.getMessage();
        }

        return output.toString();
    }

    public static void main(String[] args) throws Exception {

        executeCommand("{{.Cmd}}");

    }

}