# force Tilt to use the registry of k3s
default_registry('localhost:5000')

# avoid operation to be executed in a production grade cluster
allow_k8s_contexts('k3d-test')

# build: tell Tilt what images to build from which directories
docker_build('maestro-k8s', './')

# deploy: tell Tilt what YAML to deploy
k8s_yaml('config/kubernetes.yaml')
