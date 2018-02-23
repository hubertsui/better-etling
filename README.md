# MS Azure ACI/AKS ETL Demo

**Table of content**

[Foreword](#foreword)

[Build Docker Images](#build-docker-images)

[Deploy Azure Components with ARM Template](#deploy-azure-components-with-arm-template)

* [Create Azure Storage Account and File Share](#create-azure-storage-account-and-file-share)
* [GitHub Authorize](#github-authorize)
* [Deploy Azure Components](#deploy-azure-components)

[Check Result of the Demo](#check-result-of-the-demo)

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

1. Clone the repo, open with Visual Studio Code

2. Execute the commands below which builds and pushes extracting image:

   ```powershell
   cd cmd/extracting
   docker build -t YOURDOCKERACCOUNTNAME/extracting:1.1 .
   docker push YOURDOCKERACCOUNTNAME/extracting:1.1
   ```
   ![](Images/docker-01.png)

2. Execute the commands below which builds and pushes transforming image:

   ```powershell
   cd ../extracting
   docker build -t YOURDOCKERACCOUNTNAME/transforming:1.2 .
   docker push YOURDOCKERACCOUNTNAME/transforming:1.2
   ```
   ![](Images/docker-02.png)

3. Execute the commands below which builds and pushes loading image:

   ```powershell
   cd ../extracting
   docker build -t YOURDOCKERACCOUNTNAME/loading:1.2 .
   docker push YOURDOCKERACCOUNTNAME/loading:1.2
   ```
   ![](Images/docker-03.png)

4. Execute the commands below which builds and pushes rendering image:

   ```powershell
   cd ../extracting
   docker build -t YOURDOCKERACCOUNTNAME/rendering:1.1 .
   docker push YOURDOCKERACCOUNTNAME/rendering:1.1
   ```
   ![](Images/docker-04.png)

## Deploy Azure Components with ARM Template

### Create Azure Storage Account and File Share

1. Open the Shell in Azure Portal

   ![](Images/deploy-01.png)

2. Execute the commands below to choose subscription:

   ```powershell
   az account set --subscription SELECTED_SUBSCRIPTION_ID
   ```

3. Execute the commands below which creates a resource group:

   ```powershell
   ACI_PERS_RESOURCE_GROUP=MSAzure-ACIAKS-ETL-Demo
   ACI_PERS_STORAGE_ACCOUNT_NAME=mystorageaccount$RANDOM
   ACI_PERS_LOCATION=eastus
   ACI_PERS_SHARE_NAME=acishare
   az group create --location eastus --name $ACI_PERS_RESOURCE_GROUP
   ```

4. Execute the commands below which create the storage account:

   ```powershell
   az storage account create \
    --resource-group $ACI_PERS_RESOURCE_GROUP \
    --name $ACI_PERS_STORAGE_ACCOUNT_NAME \
    --location $ACI_PERS_LOCATION \
    --sku Standard_LRS
   ```

5. Execute the commands below which create the file share:

   ```powershell
   export AZURE_STORAGE_CONNECTION_STRING=`az storage account show-connection-string --resource-group $ACI_PERS_RESOURCE_GROUP --name $ACI_PERS_STORAGE_ACCOUNT_NAME --output tsv`
   az storage share create -n $ACI_PERS_SHARE_NAME
   ```

6. Execute the commands below which shows the **Storage Account Name**:

   ```powershell
   echo $ACI_PERS_STORAGE_ACCOUNT_NAME
   ```
   ![](Images/deploy-02.png)

### GitHub Authorize

1. Generate Token

   - Open [https://github.com/settings/tokens](https://github.com/settings/tokens) in your web browser.

   - Sign into your GitHub account where you forked this repository.

   - Click **Generate Token**.

   - Enter a value in the **Token description** text box.

   - Select the following s (your selections should match the screenshot below):

     - repo (all) -> repo:status, repo_deployment, public_repo
     - admin:repo_hook -> read:repo_hook

     ![](Images/github-new-personal-access-token.png)

   - Click **Generate token**.

   - Copy the token.

2. Add the GitHub Token to Azure in the Azure Resource Explorer

   - Open [https://resources.azure.com/providers/Microsoft.Web/sourcecontrols/GitHub](https://resources.azure.com/providers/Microsoft.Web/sourcecontrols/GitHub) in your web browser.

   - Log in with your Azure account.

   - Selected the correct Azure subscription.

   - Select **Read/Write** mode.

   - Click **Edit**.

   - Paste the token into the **token parameter**.

     ![](Images/update-github-token-in-azure-resource-explorer.png)

   - Click **PUT**.

### Deploy Azure Components

1. Fork this repository to your GitHub account.

2. Click the Deploy to Azure Button:

   [![Deploy to Azure](https://camo.githubusercontent.com/9285dd3998997a0835869065bb15e5d500475034/687474703a2f2f617a7572656465706c6f792e6e65742f6465706c6f79627574746f6e2e706e67)](https://portal.azure.com/#create/Microsoft.Template/uri/https%3A%2F%2Fraw.githubusercontent.com%2Fhubertsui%2Fbetter-etling%2Fmaster%2Fazuredeploy.json)

3. Fill in the values on the deployment page:

   ![](Images/deploy-03.png)

   You have collected most of the values in previous steps. For the rest parameters:

   * **Storage Account Name**: the name of the storage account you just created
   * **Storage Share Name**: the name of the file share. In this case, it's acishare. 
   * **Administrator Login**:  the user name of the postgres database. The default value is postgres
   * **Administrator Login Password**: the password of the postgres database. The default value is password123!@#. It must contain numbers, letters and symbols.
   * **Extracting Container Image**: the docker image you build for extracting container. The default value is hubertsui/extracting:1.1.
   * **Transforming Container Image**: the docker image you build for transforming container. The default value is hubertsui/transforming:1.2.
   * **Loading Container Image**: the docker image you build for loading container. The default value is hubertsui/loading:1.2.
   * **Rendering Container Image**: the docker image you build for rendering container. The default value is hubertsui/rendering:1.1.
   * Check **I agree to the terms and conditions stated above**.

4. Click **Purchase**.

## Chek Result of the Demo

1. Open the Resource Group **MSAzure-ACIAKS-ETL-Demo** in Auzre Portal

   ![](Images/deploy-04.png)

2. Click the MS-ACIAKS-ETLContainerGroups, statuses of containers should be like this:

   ![](Images/deploy-05.png)

3. Open the IP address of containers, and the web page should be like this:

   ![](Images/deploy-06.png)

## References
1. Akari Asai, Sara Evensen, Behzad Golshan, Alon Halevy, Vivian Li, Andrei Lopatenko, Daniela Stepanov, Yoshihiko Suhara, Wang-Chiew Tan, Yinzhan Xu, 
``HappyDB: A Corpus of 100,000 Crowdsourced Happy Moments'', LREC '18, May 2018. (to appear)