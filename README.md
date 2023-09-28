# Receipt Processor

 This processor reads the retail and purchase information from receipts and generate unique ID for each receipt and calculate reward points for each receipt.

 I have included a few valid examples as well as some edge cases to test the processor's functionality.

## Environment Setup and Application Run
This project is written in Go. To run this project, you need to have Go installed and set up on your local machine or IDE (I used IntelliJ). 

In your local IDE, search for the "Go" plugin and install it.

On the GitHub website, you can either download or clone the project and open up in your local IDE.

Once the Go is installed and project is downloaded/cloned, follow these steps to run the application:

1. **Navigate to Repository**: Open your terminal and navigate to the repository folder.


2. **Initialize Go Modules (if needed):** If you're using Go modules, ensure that you've initialized your project by running the following command:
   ``` bash
   go mod init your_project_name
   ```
      This command initializes a go.mod file to manage your project's dependencies.

3. **Install Dependencies**: Use the following command to download and install the project's dependencies:
   ``` bash
   go mod tidy
   ```

4. **Run Application**: Execute the following command to run the application:
   ``` bash
   go run main.go
   ```

## Application Test
### Unit Testing
To maintain the integrity and correctness of our application, I've included a thorough set of unit tests. These tests help ensure that the software components behave as expected. 
If you'd like to run these tests and check their results, follow the instructions below:

Open up the Terminal within IDE and initiate the unit test by entering:
   ``` bash
   go test main_test.go
   ```

To see the details of each test case, such as if the test case is passed or certain message is printed, please command:
   ``` bash
   go test -v main_test.go
   ```

### Integration Testing
To perform integration testing, you can use Postman or API tester (I have used Talend). 

To process the receipt, ensure to select POST as HTTP method and http://localhost:8080/receipts/process, then in the Body section, enter a test case from main_test.go or design your own.
After sent, you should receive the receipt ID. Copy the ID for retrieving points purpose. 

To retrieve reward points, ensure to select GET method and http://localhost:8080/receipts/{id}/points, be sure to replace the {id} with the ID that you copied.
After sent, you should receive either 200 or 400 and the response body. 


