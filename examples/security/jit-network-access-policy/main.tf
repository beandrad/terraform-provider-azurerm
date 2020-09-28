provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

# create vm
resource "azurerm_virtual_network" "vnet" {
  name                = "${var.prefix}network"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
}

resource "azurerm_subnet" "subnet" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.rg.name
  virtual_network_name = azurerm_virtual_network.vnet.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "inet" {
  name                = "${var.prefix}-nic"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.subnet.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "main" {
  name                             = "${var.prefix}-vm2"
  location                         = azurerm_resource_group.rg.location
  resource_group_name              = azurerm_resource_group.rg.name
  network_interface_ids            = [azurerm_network_interface.inet.id]
  vm_size                          = "Standard_A1_v2"
  delete_os_disk_on_termination    = true
  delete_data_disks_on_termination = true

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
  storage_os_disk {
    name              = "disk1"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }
  os_profile {
    computer_name  = "hostname"
    admin_username = "isabel"
    admin_password = "12345&Isabel"
  }
  os_profile_linux_config {
    disable_password_authentication = false
  }
}

resource "azurerm_jit_network_access_policies" "jit" {
  asc_location = "${var.location}"
  resource_group_name = azurerm_resource_group.rg.name
  kind = "Basic"
  virtual_machines {
      id = azurerm_virtual_machine.main.id
      ports {
          port = 22
          allowed_source_address_prefix = "*"
          max_request_access_duration = "PT3H"
      }
      
  }
  
}

