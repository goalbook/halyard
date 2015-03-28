FROM google/golang:latest
ADD halyard halyard
CMD ["./halyard"]



