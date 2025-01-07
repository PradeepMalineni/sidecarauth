import torch
import torch.nn as nn
import torch.optim as optim
import matplotlib.pyplot as plt

# Step 1: Generate sample data
x = torch.tensor([[1.0], [2.0], [3.0], [4.0]])  # Input feature
y = torch.tensor([[2.0], [4.0], [6.0], [8.0]])  # Target output

# Reshape data for LSTM (sequence_length, batch_size, input_dim)
x = x.view(-1, 1, 1)
y = y.view(-1, 1, 1)

# Plot initial data
plt.figure(figsize=(8, 6))
plt.scatter(x.view(-1).numpy(), y.view(-1).numpy(), color='blue', label='Original Data')
plt.xlabel('Input (x)')
plt.ylabel('Output (y)')
plt.title('Input vs Output Data')
plt.legend()
plt.show()

# Step 2: Define an LSTM-based model
class LSTMModel(nn.Module):
    def __init__(self):
        super(LSTMModel, self).__init__()
        self.lstm = nn.LSTM(input_size=1, hidden_size=50, num_layers=1, batch_first=True)
        self.fc = nn.Linear(50, 1)

    def forward(self, x):
        out, _ = self.lstm(x)  # LSTM layer
        out = self.fc(out)    # Fully connected layer
        return out

model = LSTMModel()

# Step 3: Define loss function and optimizer
criterion = nn.MSELoss()
optimizer = optim.SGD(model.parameters(), lr=0.01)

# Step 4: Training loop
iterations = 1000
for i in range(iterations):
    # Forward pass
    y_pred = model(x)

    # Compute loss
    loss = criterion(y_pred, y)

    # Backward pass and optimization
    optimizer.zero_grad()
    loss.backward()
    optimizer.step()

    # Print loss every 100 iterations
    if (i + 1) % 100 == 0:
        print(f'Iteration {i+1}, Loss: {loss.item()}')

# Final parameters
for name, param in model.named_parameters():
    print(f'{name}: {param.data}')

# Plot final predictions
with torch.no_grad():
    y_final = model(x)
    plt.figure(figsize=(8, 6))
    plt.scatter(x.view(-1).numpy(), y.view(-1).numpy(), color='blue', label='Original Data')
    plt.plot(x.view(-1).numpy(), y_final.view(-1).numpy(), color='red', label='Fitted Line')
    plt.xlabel('Input (x)')
    plt.ylabel('Output (y)')
    plt.title('Final Model Fit with LSTM')
    plt.legend()
    plt.show()
