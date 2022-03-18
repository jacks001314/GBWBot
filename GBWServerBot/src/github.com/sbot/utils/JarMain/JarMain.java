public class JarMain {

    public JarMain() {
    }

    private static String executeCommand(String cmd) {

        StringBuilder sb = new StringBuilder();

        try {

            ProcessBuilder processBuilder = new ProcessBuilder();

            processBuilder.command(cmd.split(","));

            Process process = processBuilder.start();

            process.waitFor();

        } catch (Exception e) {
            return e.getMessage();
        }

        return sb.toString();
    }

    public static void main(String[] args) throws Exception {

        executeCommand("whoami");

    }
}
