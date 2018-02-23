# MS Azure ACI/AKS ETL Demo

**Table of content**

[Foreword](#foreword)

[Build Docker Images](#build-docker-images)

[Create Azure Storage Account and File Share](#create-azure-storage-account-and-file-share)

[Deploy Azure Components](#deploy-azure-components)

[Check the Demo](#check-result-of-the-demo)

[References](#references)

## Foreword

Making ETLing simpler using containers as a plug and play model.

The application with separate common pieces of the ETLing process into separate docker containers. For example, the unzip container in the project will takes a link as input then downloads and unzips the file at that link. A separate container can take a csv location as input and put it into a postgres database. This allows for plug and play ETLing pipelines for the data.

Using ACI, a user can define container groups with the exact elements they want. For example put the unzip and postgres modules together and can download a zip file from a datasource, unzip it then feed it into a databases all without writing a line of code. Also only pay per second usign the ACI instance. 

This document will guide you to deploy the solution to your environment.

First, an Azure AAD is required to register the app registrations. In this document, the Azure AAD will be called "ETL AAD", and an account in ETL AAD will be called ETL work account.

* All app registrations should be created in the ETL AAD. 

An Azure Subscription is required to deploy the Azure components. We will use the [ARM Template](azuredeploy.json) to deploy these Azure components automatically. 

If you want to build your own docker images, please clone this repo, and install docker on your own computer.

## Build Docker Images

1. Clone the repo with Visual Studio Code.

2. Open **View > Integrated Terminal**.

3. Execute the commands below which builds and pushes extracting image to Docker Hub.

   ```powershell
   cd cmd/extracting
   docker build -t YOURDOCKERACCOUNTNAME/extracting:1.0 .
   docker push YOURDOCKERACCOUNTNAME/extracting:1.0
   ```

4. Execute the commands below which builds and pushes transforming image to Docker Hub.

   ```powershell
   cd ../extracting
   docker build -t YOURDOCKERACCOUNTNAME/transforming:1.0 .
   docker push YOURDOCKERACCOUNTNAME/transforming:1.0
   ```

5. Execute the commands below which builds and pushes loading image to Docker Hub.

   ```powershell
   cd ../extracting
   docker build -t YOURDOCKERACCOUNTNAME/loading:1.0 .
   docker push YOURDOCKERACCOUNTNAME/loading:1.0
   ```

6. Execute the commands below which builds and pushes rendering image to Docker Hub.

   ```powershell
   cd ../extracting
   docker build -t YOURDOCKERACCOUNTNAME/rendering:1.0 .
   docker push YOURDOCKERACCOUNTNAME/rendering:1.0
   ```

## Create Azure Storage Account and File Share

1. Open the Shell in Azure Portal

   ![](Images/deploy-01.png)

2. Execute the commands below to choose subscription:

   ```powershell
   az account set --subscription SELECTED_SUBSCRIPTION_ID
   ```

3. Execute the commands below which creates a new resource group:

   > Note: Change the placeholder `[RESOURCE_GROUP_NAME]` to a new resource group to be created.
   
   ```powershell
   ACI_PERS_RESOURCE_GROUP=[RESOURCE_GROUP_NAME]
   ACI_PERS_STORAGE_ACCOUNT_NAME=mystorageaccount$RANDOM
   ACI_PERS_LOCATION=eastus
   ACI_PERS_SHARE_NAME=acishare
   az group create --location eastus --name $ACI_PERS_RESOURCE_GROUP
   ```

4. Execute the commands below which creates the storage account:

   ```powershell
   az storage account create \
    --resource-group $ACI_PERS_RESOURCE_GROUP \
    --name $ACI_PERS_STORAGE_ACCOUNT_NAME \
    --location $ACI_PERS_LOCATION \
    --sku Standard_LRS
   ```

5. Execute the commands below which creates the file share:

   ```powershell
   export AZURE_STORAGE_CONNECTION_STRING=`az storage account show-connection-string --resource-group $ACI_PERS_RESOURCE_GROUP --name $ACI_PERS_STORAGE_ACCOUNT_NAME --output tsv`
   az storage share create -n $ACI_PERS_SHARE_NAME
   ```

6. Execute the commands below which shows the storage account created previously:

   ```powershell
   echo $ACI_PERS_STORAGE_ACCOUNT_NAME
   ```
   ![](Images/deploy-02.png)

## Deploy Azure Components

1. Click this button to navigate to Azure portal deployment page.

   [![Deploy to Azure](https://azuredeploy.net/deploybutton.png)](https://portal.azure.com/#create/Microsoft.Template/uri/https%3A%2F%2Fraw.githubusercontent.com%2Fhubertsui%2Fbetter-etling%2Fmaster%2Fazuredeploy.json)

2. Fill in the values on the deployment page:
   * **Storage Account Name**: the name of the storage account you just created.
   * **Storage Share Name**: the name of the file share. In this case, it's `acishare` .
   * **Administrator Login**:  the user name of the postgres database.
   * **Administrator Login Password**: the password of the postgres database, it must meet the complexity requirements, e.g. `password123!@#` .
   * **Extracting Container Image**: the docker image you build for extracting container.
   * **Transforming Container Image**: the docker image you build for transforming container.
   * **Loading Container Image**: the docker image you build for loading container.
   * **Rendering Container Image**: the docker image you build for rendering container.
   * Check **I agree to the terms and conditions stated above**.

   ![](Images/deploy-03.png)

3. Click **Purchase**.

## Chek the Demo

1. Open the resource group just created.

   ![](Images/deploy-04.png)

2. Click the container group **MS-ACIAKS-ETLContainerGroups**, state of the first 3 containers are **Terminated**, while the last one is always **Running** to serve the Word Cloud.

   ![](Images/deploy-05.png)

3. Copy the **IP address** from the container group blade, and open it in the browser to check the Word Cloud.

   ![](Images/deploy-06.png)

## References
1. Akari Asai, Sara Evensen, Behzad Golshan, Alon Halevy, Vivian Li, Andrei Lopatenko, Daniela Stepanov, Yoshihiko Suhara, Wang-Chiew Tan, Yinzhan Xu, 
``HappyDB: A Corpus of 100,000 Crowdsourced Happy Moments'', LREC '18, May 2018. (to appear)
